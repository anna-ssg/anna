#This script can be utilised to clean the test output data

cd test/
rm -rf `find . -type d -name rendered`
cd ../