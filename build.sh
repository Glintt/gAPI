# create build folder
rm -rf bin
mkdir -p bin/api/configs bin/api/src bin/dashboard

# build API
echo "######    Building API    #######"
cd api
echo "1. Building ...."
sh build.sh
echo "2. Copying sources"
cp -r . ../bin/api/src
echo "3. Adding config folder"
cp -r configs-example/* ../bin/api/configs

# build dashboard
#echo "######    Building Dashboard    #######"
#cd ../dashboard
#echo "1. Installing dependencies ...."
#npm install
#echo "   Dependencies installation finished!"
#echo "2. Building dashboard ...."
#npm run build
#echo "   Finished to build dashboard!"


# copy generated sources to final build folder
echo "######    Finalizing build    #######"
cd ..

echo "1. Copying API executables"
cp -r api/bin/* bin/api/bin
echo "2. Copying dashboard distribution"
cp -r dashboard/public/ bin/dashboard/public
cp -r dashboard/src/ bin/dashboard/src
cp dashboard/package.json bin/dashboard/package.json
cp dashboard/index.js bin/dashboard/index.js

# remove not needed files
echo "3. Removing temp files"
rm -rf api/bin
rm -rf dashboard/dist

echo "4. Add docker"
cp docker-compose.yml bin/docker-compose.yml
cp .env.example bin/.env.example

cp dashboard/.eslintrc.js bin/dashboard/.eslintrc.js
cp dashboard/.postcssrc.js bin/dashboard/.postcssrc.js
cp dashboard/babel.config.js bin/dashboard/babel.config.js

cp dashboard/Dockerfile bin/dashboard/Dockerfile
cp dashboard/.dockerignore bin/dashboard/.dockerignore

echo "#####################################"
echo "######      ENDED BUILD       #######"
echo "#####################################"
