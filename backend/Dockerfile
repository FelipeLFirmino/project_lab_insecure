# Estágio 1: Build (Compilação)
# CORREÇÃO: Mudamos de 1.21 para 1.24
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copia os arquivos de dependência e baixa elas
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código-fonte
COPY . .

# Compila o binário final. 
RUN CGO_ENABLED=0 GOOS=linux go build -o /easytrip-backend .

# Estágio 2: Run (Execução)
FROM alpine:latest

# Copia o binário compilado do estágio de build
COPY --from=build /easytrip-backend /easytrip-backend

# Expõe a porta que seu backend Go usa
EXPOSE 8080

# O comando para rodar a aplicação (o binário)
CMD ["/easytrip-backend"]