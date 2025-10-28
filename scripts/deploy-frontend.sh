#!/bin/bash

# Script para deploy do frontend para S3 + CloudFront
# Uso: ./scripts/deploy-frontend.sh [bucket-name] [api-url]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

# Verificar parâmetros
if [ $# -lt 1 ]; then
    error "Uso: $0 <bucket-name> [api-url]"
    echo "Exemplo: $0 financial-control-homolog-frontend https://api.example.com/api/v1"
    exit 1
fi

BUCKET_NAME=$1
API_URL=${2:-"https://api.example.com/api/v1"}

log "Iniciando deploy do frontend..."
log "Bucket: $BUCKET_NAME"
log "API URL: $API_URL"

# Verificar se AWS CLI está instalado
if ! command -v aws &> /dev/null; then
    error "AWS CLI não está instalado. Instale com: brew install awscli"
fi

# Verificar se está logado no AWS
if ! aws sts get-caller-identity &> /dev/null; then
    error "Não está logado no AWS. Execute: aws configure"
fi

# Verificar se o bucket existe
if ! aws s3 ls "s3://$BUCKET_NAME" &> /dev/null; then
    error "Bucket $BUCKET_NAME não existe. Crie primeiro com Terraform."
fi

# Navegar para o diretório do frontend
cd "$(dirname "$0")/../src/frontend"

log "Instalando dependências..."
npm ci

log "Configurando variáveis de ambiente..."
export VITE_API_URL="$API_URL"

log "Fazendo build do frontend..."
npm run build

if [ ! -d "dist" ]; then
    error "Build falhou - diretório dist não foi criado"
fi

log "Fazendo upload para S3..."
aws s3 sync dist/ "s3://$BUCKET_NAME" --delete

log "Invalidando cache do CloudFront..."
# Buscar o distribution ID do bucket
DISTRIBUTION_ID=$(aws cloudfront list-distributions --query "DistributionList.Items[?Origins.Items[0].DomainName=='$BUCKET_NAME.s3-website-us-east-1.amazonaws.com'].Id" --output text)

if [ -z "$DISTRIBUTION_ID" ] || [ "$DISTRIBUTION_ID" = "None" ]; then
    warn "Não foi possível encontrar a distribuição CloudFront automaticamente"
    warn "Você pode precisar invalidar o cache manualmente no console AWS"
else
    log "Invalidando cache da distribuição: $DISTRIBUTION_ID"
    aws cloudfront create-invalidation --distribution-id "$DISTRIBUTION_ID" --paths "/*"
fi

log "Deploy concluído com sucesso!"
log "Frontend disponível em: https://$BUCKET_NAME.s3-website-us-east-1.amazonaws.com"
log "Ou via CloudFront (se configurado): https://$DISTRIBUTION_ID.cloudfront.net"
