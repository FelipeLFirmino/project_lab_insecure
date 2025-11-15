# Estágio 1: Build (Compilação do Angular)
# Usa node:22-alpine para garantir compatibilidade com as versões recentes do Angular
FROM node:22-alpine AS build

WORKDIR /app

# Copia os arquivos de dependência e instala
COPY package*.json ./
RUN npm install

# Copia todo o resto do código-fonte
COPY . .

# Roda o build de produção do Angular. 
# O output é geralmente em /app/dist/easytrip/browser/
RUN npm run build -- --configuration production

# Estágio 2: Serve (Servidor Nginx)
# Usamos uma imagem Nginx leve para servir os arquivos estáticos
FROM nginx:1.25-alpine

# 1. Copia o arquivo de configuração customizado do Nginx.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 2. COPIA FINAL CORRIGIDA: Copia o conteúdo da subpasta 'browser' 
COPY --from=build --chown=nginx:nginx /app/dist/easytrip/browser /usr/share/nginx/html


# O Nginx escuta na porta 80 por padrão
EXPOSE 80