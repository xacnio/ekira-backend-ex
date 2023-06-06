DOCKER_FOLDER=".docker"
PLATFORM_FOLDER="platform/database/docker-entrypoint-initdb"
rm -rf ekira-backend/
git clone --depth 1 https://github.com/xacnio/ekira-backend-ex.git
cp .env.prod ekira-backend/.env.prod
cd ekira-backend/
cp $DOCKER_FOLDER/nginx.conf $DOCKER_FOLDER/Dockerfile $DOCKER_FOLDER/docker-compose.yml .
cp -f $DOCKER_FOLDER/docker-compose.yml $DOCKER_FOLDER/nginx.conf VERSION ..
mkdir -p ../$PLATFORM_FOLDER
cp -f $PLATFORM_FOLDER/*.sql ../$PLATFORM_FOLDER/
sh scripts/build.sh
cd .. && rm -rf ekira-backend
docker-compose --env-file .env.prod up -d
