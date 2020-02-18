# test sending image
curl http://localhost:3000/upload-image -F "file=@test.jpg" -vvv

# test sending zip
curl http://localhost:3000/upload-zip -F "file=@test.zip" -vvv

# test generating thumbnails
curl --location --request POST 'localhost:3000/generate-thumbnails' \
--header 'Content-Type: application/json' \
--data-raw '{
	"imagePath": "images/ecd1c63a-ec5f-499b-bfb4-8479f6d31b1a.jpg",
	"widthSize": 64
}'