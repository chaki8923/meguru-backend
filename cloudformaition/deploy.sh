#!/bin/bash

# deploy.sh - CloudFormation自動デプロイスクリプト
set -e

# 色付きログ用の定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# デフォルト値
STACK_NAME="meguru-stack"
REGION="ap-northeast-1"
TEMPLATE_FILE="template.yaml"
PARAMETERS_FILE="parameters.json"
TAGS_FILE="tags.json"

# ヘルプ表示
show_help() {
    cat << EOF
AWS CloudFormation自動デプロイスクリプト

使用方法:
    $0 [OPTIONS]

オプション:
    -s, --stack-name NAME     スタック名 (デフォルト: meguru-stack)
    -r, --region REGION       AWSリージョン (デフォルト: ap-northeast-1)
    -t, --template FILE       CloudFormationテンプレートファイル
    -p, --parameters FILE     パラメータファイル
    -d, --delete             スタックを削除
    -h, --help               このヘルプを表示

例:
    # 新規デプロイ
    $0 -s my-app -t template.yaml -p parameters.json

    # スタック削除
    $0 -s my-app --delete
EOF
}

# コマンドライン引数解析
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--stack-name)
            STACK_NAME="$2"
            shift 2
            ;;
        -r|--region)
            REGION="$2"
            shift 2
            ;;
        -t|--template)
            TEMPLATE_FILE="$2"
            shift 2
            ;;
        -p|--parameters)
            PARAMETERS_FILE="$2"
            shift 2
            ;;
        -d|--delete)
            DELETE_STACK=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "未知のオプション: $1"
            show_help
            exit 1
            ;;
    esac
done

# AWS CLI確認
check_aws_cli() {
    log_info "AWS CLI確認中..."
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLIがインストールされていません"
        exit 1
    fi
    
    # AWS認証確認
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS認証が設定されていません"
        exit 1
    fi
    
    log_success "AWS CLI設定OK"
}

# ファイル存在確認
check_files() {
    if [[ ! -f "$TEMPLATE_FILE" ]]; then
        log_error "テンプレートファイルが見つかりません: $TEMPLATE_FILE"
        exit 1
    fi
    
    if [[ ! -f "$PARAMETERS_FILE" ]]; then
        log_error "パラメータファイルが見つかりません: $PARAMETERS_FILE"
        exit 1
    fi
    
    log_success "必要ファイル確認完了"
}

# スタック存在確認
check_stack_exists() {
    aws cloudformation describe-stacks \
        --stack-name "$STACK_NAME" \
        --region "$REGION" &> /dev/null
}

# テンプレート検証
validate_template() {
    log_info "CloudFormationテンプレート検証中..."
    
    if aws cloudformation validate-template \
        --template-body "file://$TEMPLATE_FILE" \
        --region "$REGION" &> /dev/null; then
        log_success "テンプレート検証OK"
    else
        log_error "テンプレート検証失敗"
        exit 1
    fi
}

# Change Set作成・確認
create_change_set() {
    local change_set_name="deploy-$(date +%Y%m%d-%H%M%S)"
    
    log_info "Change Set作成中: $change_set_name"
    
    if check_stack_exists; then
        local change_set_type="UPDATE"
    else
        local change_set_type="CREATE"
    fi
    
    aws cloudformation create-change-set \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --template-body "file://$TEMPLATE_FILE" \
        --parameters "file://$PARAMETERS_FILE" \
        --tags "file://$TAGS_FILE" \
        --capabilities CAPABILITY_NAMED_IAM \
        --change-set-name "$change_set_name" \
        --change-set-type "$change_set_type" > /dev/null
    
    # Change Set準備完了まで待機
    log_info "Change Set準備完了まで待機中..."
    aws cloudformation wait change-set-create-complete \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --change-set-name "$change_set_name"
    
    # Change Set内容表示
    log_info "Change Set内容:"
    aws cloudformation describe-change-set \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --change-set-name "$change_set_name" \
        --query 'Changes[*].[Action,ResourceChange.ResourceType,ResourceChange.LogicalResourceId]' \
        --output table
    
    # 実行確認
    echo
    read -p "Change Setを実行しますか? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_warning "デプロイをキャンセルしました"
        # Change Set削除
        aws cloudformation delete-change-set \
            --stack-name "$STACK_NAME" \
            --region "$REGION" \
            --change-set-name "$change_set_name" > /dev/null
        exit 0
    fi
    
    # Change Set実行
    log_info "Change Set実行中..."
    aws cloudformation execute-change-set \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --change-set-name "$change_set_name" > /dev/null
}

# デプロイ進行状況監視
monitor_deployment() {
    log_info "デプロイ進行状況監視中..."
    
    # デプロイ完了まで待機
    local start_time=$(date +%s)
    
    if check_stack_exists; then
        aws cloudformation wait stack-update-complete \
            --stack-name "$STACK_NAME" \
            --region "$REGION" &
    else
        aws cloudformation wait stack-create-complete \
            --stack-name "$STACK_NAME" \
            --region "$REGION" &
    fi
    
    local wait_pid=$!
    
    # 進行状況を定期表示
    while kill -0 $wait_pid 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        printf "\r⏱️  経過時間: %02d:%02d " $((elapsed/60)) $((elapsed%60))
        
        # 最新のイベントを取得
        local latest_event=$(aws cloudformation describe-stack-events \
            --stack-name "$STACK_NAME" \
            --region "$REGION" \
            --query 'StackEvents[0].[ResourceType,ResourceStatus]' \
            --output text 2>/dev/null || echo "Unknown Unknown")
        
        printf "| 最新: %s" "$latest_event"
        
        sleep 5
    done
    
    echo # 改行
    
    # 結果確認
    wait $wait_pid
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        local end_time=$(date +%s)
        local total_time=$((end_time - start_time))
        log_success "デプロイ完了！ (総時間: $((total_time/60))分$((total_time%60))秒)"
    else
        log_error "デプロイ失敗"
        show_error_events
        exit 1
    fi
}

# エラーイベント表示
show_error_events() {
    log_error "失敗したリソース:"
    aws cloudformation describe-stack-events \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --query 'StackEvents[?ResourceStatus==`CREATE_FAILED` || ResourceStatus==`UPDATE_FAILED`].[Timestamp,ResourceType,ResourceStatusReason]' \
        --output table
}

# 出力値表示
show_outputs() {
    log_info "スタック出力値:"
    aws cloudformation describe-stacks \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --query 'Stacks[0].Outputs' \
        --output table
    
    # App Runner URLを特別に表示
    local app_runner_url=$(aws cloudformation describe-stacks \
        --stack-name "$STACK_NAME" \
        --region "$REGION" \
        --query 'Stacks[0].Outputs[?OutputKey==`AppRunnerServiceUrl`].OutputValue' \
        --output text 2>/dev/null)
    
    if [[ -n "$app_runner_url" ]]; then
        echo
        log_success "🚀 アプリケーションURL: $app_runner_url"
        log_info "ヘルスチェック: curl $app_runner_url/health"
    fi
}

# スタック削除
delete_stack() {
    log_warning "スタック削除: $STACK_NAME"
    
    read -p "本当にスタックを削除しますか? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "削除をキャンセルしました"
        exit 0
    fi
    
    log_info "スタック削除中..."
    aws cloudformation delete-stack \
        --stack-name "$STACK_NAME" \
        --region "$REGION"
    
    log_info "削除完了まで待機中..."
    aws cloudformation wait stack-delete-complete \
        --stack-name "$STACK_NAME" \
        --region "$REGION"
    
    log_success "スタック削除完了"
}

# メイン処理
main() {
    log_info "CloudFormation自動デプロイ開始"
    log_info "スタック名: $STACK_NAME"
    log_info "リージョン: $REGION"
    
    check_aws_cli
    
    if [[ "${DELETE_STACK:-false}" == "true" ]]; then
        delete_stack
        exit 0
    fi
    
    check_files
    validate_template
    create_change_set
    monitor_deployment
    show_outputs
    
    log_success "🎉 デプロイ完了！"
}

# エラーハンドリング
trap 'log_error "スクリプト実行中にエラーが発生しました"; exit 1' ERR

# メイン処理実行
main "$@"