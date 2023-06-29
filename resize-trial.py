import boto3
import os
from PIL import Image

def resize_image(file_path, target_size):
    #image = Image.convert('RGB')
    image = Image.open(file_path)
    image = image.convert('RGB')
    image.thumbnail(target_size)
    resized_path = file_path.replace('original_folder', 'resized_folder')
    image.save(resized_path)
    return resized_path

def resize_images_in_s3_bucket(bucket_name, original_folder, resized_folder, target_size):
    s3 = boto3.client('s3')
    response = s3.list_objects_v2(Bucket=bucket_name, Delimiter="/", Prefix=original_folder)

    for obj in response['Contents']:
        file_key = obj['Key']
        file_name = os.path.basename(file_key)
        if file_name.endswith('.jpg') or file_name.endswith('.jpeg') or file_name.endswith('.png'):
            # Download the file locally
            local_file_path = '/tmp/' + file_name
            s3.download_file(bucket_name, file_key, local_file_path)

            # Resize the image
            resized_path = resize_image(local_file_path, target_size)

            # Upload the resized image to S3
            resized_key = file_key.replace(original_folder, resized_folder)
            s3.upload_file(resized_path, bucket_name, resized_key)

            # Remove the local resized image
            os.remove(resized_path)

            print(f"Resized and uploaded {file_name} to {resized_key}")

    print("All images resized and uploaded successfully.")

# Set your AWS credentials - Auth very important
boto3.setup_default_session(aws_access_key_id='ACCESS_ID',
                           aws_secret_access_key='SECRET_ACCESS_ID',
                            region_name='ap-south-1')

# Specify your bucket name, original and resized folder paths, and target size
bucket_name = 'bucket name'
original_folder = 'source folder within the bucket'
resized_folder = 'desstination folder without the bucket'
target_size = (512, 512)  # Specify your desired target size here

# Call the function to resize images in the S3 bucket
resize_images_in_s3_bucket(bucket_name, original_folder, resized_folder, target_size)
