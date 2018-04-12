# create build folder
mkdir -p bin/api/configs bin/dashboard/dist

# build API
echo "######    Building API    #######"
cd api
echo "1. Building ...."
sh build.sh
echo "2. Adding config folder"
cp -r configs-example/* ../bin/api/configs

# build dashboard
echo "######    Building Dashboard    #######"
cd ../dashboard
echo "1. Installing dependencies ...."
npm install
echo "   Dependencies installation finished!"
echo "2. Building dashboard ...."
npm run build
echo "   Finished to build dashboard!"


# copy generated sources to final build folder
echo "######    Finalizing build    #######"
cd ..

echo "1. Copying API executables"
cp -r api/bin/* bin/api
echo "2. Copying dashboard distribution"
cp -r dashboard/dist/ bin/dashboard
cp dashboard/package.json bin/dashboard/package.json
cp dashboard/index.js bin/dashboard/index.js

# remove not needed files
echo "3. Removing temp files"
rm -rf api/bin
rm -rf dashboard/dist

echo "#####################################"
echo "######      ENDED BUILD       #######"
echo "#####################################"
