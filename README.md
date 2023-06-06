# ekira-backend

## Used Technologies
- [Golang](https://golang.org/)
- [Gorm](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Nginx](https://www.nginx.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Redis](https://redis.io/)
- [Cloudflare](https://www.cloudflare.com/)
- [VerifyKit](https://verifykit.com/)
- [Stripe](https://stripe.com/)

# Docker Installation
- [Docker](https://www.docker.com/)

## Configuration
### Project Configuration
```bash
sudo cp .env.example .env.prod
nano .env.prod
```
### Run
```bash
docker build -t xacnio/ekira-backend .
docker-compose --env-file .env.prod up -d
```

# Manual Installation

## Requirements
- Ubuntu 20.04.5 LTS
- Go 1.18 or higher
- PostgreSQL 15.1 or higher
- ImageMagick 6.9.10.23
- Nodejs 19 - PM2 5.2.2 (production)

## Installation
### Golang Install
```bash
wget https://go.dev/dl/go1.19.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
```

### PostgreSQL Install
```bash
sudo apt-get update
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get -y install postgresql
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### ImageMagick Install
```bash
sudo apt-get update
sudo apt-get install imagemagick libmagickwand-dev pkg-config
sudo apt install build-essential make gcc
sudo apt install libx11-dev libxext-dev zlib1g-dev libpng-dev libjpeg-dev libfreetype6-dev libxml2-dev
wget https://download.imagemagick.org/archive/ImageMagick-7.1.1-11.tar.gz
tar -xvf ImageMagick-7.1.1-11.tar.gz
cd ImageMagick-7.1.1-11/
./configure
make -j4
sudo make install
sudo ldconfig /usr/local/lib
magick -version
```

### Nginx Install
```bash
sudo apt-get update
sudo apt-get install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### Nodejs, NVM, PM2 ve Other Packages Install
```bash
curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash
source ~/.bashrc
nvm install 19.0.0
nvm use 19.0.0
sudo npm install pm2@latest -g
sudo apt install make -y
```

## Configuration
### Nginx Configuration
```bash
sudo nano /etc/nginx/conf.d/ekira.nginx.conf
```

```nginx
server {
    listen      80;
    server_name api.e-kira.tk www.e-kira.tk;
    error_log  /var/log/nginx/api.e-kira.tk.error.log;
    access_log  /var/log/nginx/api.e-kira.tk.access.log;

    location / {
        proxy_no_cache 1;
        proxy_cache_bypass 1;
        proxy_set_header   Host     $host;
        client_max_body_size 10M;
        proxy_pass      http://127.0.0.1:5000;
    }

    location ~ /\.ht    {return 404;}
    location ~ /\.svn/  {return 404;}
    location ~ /\.git/  {return 404;}
    location ~ /\.hg/   {return 404;}
    location ~ /\.bzr/  {return 404;}
}
```

```bash
sudo service nginx restart
```

### Project Setup
#### Production
```bash
echo "export RUN_TYPE=PROD" >> ~/.bashrc
source ~/.bashrc
```

#### Development
```bash
echo "export RUN_TYPE=DEV" >> ~/.bashrc
source ~/.bashrc
```

#### Project Configuration
##### Production
```bash
sudo cp .env.example .env.prod
nano .env.prod
```

##### Development
```bash
sudo cp .env.example .env.dev
nano .env.dev
```

```
# Server settings:
API_URL="https://api.e-kira.tk"
SERVER_HOST="127.0.0.1:5000"
SERVER_URL="api.e-kira.tk"
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY="secret"
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=518400

# Database settings:
DB_SERVER_URL="postgres://postgres:secret-password@localhost:5432?sslmode=disable"
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=2
```

### Build
```bash
make swag
make build
```

### Run
```bash
pm2 start --name ekira-backend ./build/ekira-backend
```
or
```bash
make run
```

### Docker
Not available yet.
