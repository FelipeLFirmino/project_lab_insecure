
-----

````markdown
# EasyTrip: API e Frontend (Full Stack Dockerizado)

Este é o projeto **EasyTrip**, uma aplicação para organização colaborativa de viagens em grupo, composta por:

1.  **Backend API:** Desenvolvido em **Go (Golang)**, utilizando uma arquitetura em camadas (`handlers`, `services`, `repositories`) e **PostgreSQL** como banco de dados.
2.  **Frontend:** Aplicação Single Page Application (SPA) desenvolvida em **Angular**, servida por **Nginx**.

---

## 🚀 Como Rodar a Aplicação (Docker Compose)

O projeto está totalmente configurado para ser executado usando **Docker Compose**, o que elimina a necessidade de instalar Go, Node.js ou PostgreSQL diretamente em sua máquina.

### Pré-requisitos

Certifique-se de que as seguintes ferramentas estão instaladas em seu sistema:

-   **Docker:** [Download e Instalação Oficial](https://docs.docker.com/get-docker/)
-   **Docker Compose:** Geralmente incluído na instalação do Docker Desktop.

### 🛠️ Configuração Inicial

#### Passo 1: Clone o Repositório

Clone o repositório para o seu ambiente local:

```bash
git clone <url-do-seu-repositorio>
cd <nome-do-seu-projeto>
````

#### Passo 2: Configure as Variáveis de Ambiente

O Docker Compose precisa de algumas variáveis para configurar o banco de dados. Crie um arquivo na **raiz do projeto** (na mesma pasta onde está o `docker-compose.yml`) chamado **`.env`**.

Defina as credenciais para o banco de dados.

```ini
# Variáveis de Ambiente do Banco de Dados (Postgres)
DB_HOST=db
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=sua_senha_aqui
DB_NAME=project_lab

# Variável de Ambiente do Backend (Go)
# A porta 8080 é o padrão do backend
SERVER_PORT=8080 
```

⚠️ **Importante:** Substitua `sua_senha_aqui` por uma senha segura. Este arquivo deve estar no `.gitignore`.

-----

### ▶️ Inicializando a Aplicação

Execute estes três comandos na raiz do projeto, em sequência, para garantir que todas as imagens sejam construídas e inicializadas corretamente:

#### 1\. Limpeza (Opcional, mas recomendado na primeira vez)

Garante que não haja contêineres antigos ou redes em conflito.

```bash
docker-compose down -v
```

#### 2\. Build (Compilação das Imagens)

Este passo compila o Backend (Go) e o Frontend (Angular/Nginx).

```bash
docker-compose build
```

#### 3\. Execução (Start)

Inicia todos os três serviços (`db`, `backend`, `frontend`) em modo *detached* (segundo plano).

```bash
docker-compose up -d
```

### ✅ Acesso à Aplicação

Todos os serviços estarão rodando e conectados:

| Serviço | Porta Exposta (Host) | URL de Acesso |
| :--- | :--- | :--- |
| **Frontend (Angular/Nginx)** | `4200` | **http://localhost:4200** |
| **Backend (Go API)** | `8080` | `http://localhost:8080/...` |
| **Banco de Dados (Postgres)** | `5432` | Acesso interno/local (Opcional) |

Acesse **`http://localhost:4200`** para interagir com o EasyTrip\!

-----

### 🛑 Parando a Aplicação

Para parar e remover todos os contêineres e a rede definidos pelo Compose, execute:

```bash
docker-compose down
```

```
```