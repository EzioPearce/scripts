import boto3
import csv

def get_s3_object_ids(bucket_name, folder_name):
    s3_client = boto3.client('s3')
    
    # List objects in the specified bucket and folder
    response = s3_client.list_objects_v2(Bucket=bucket_name, Prefix=folder_name)
    
    object_ids = []
    if 'Contents' in response:
        for item in response['Contents']:
            object_ids.append(item['Key'])
    
    return object_ids

def write_to_csv(file_path, object_ids):
    with open(file_path, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['Object ID'])
        writer.writerows([[object_id] for object_id in object_ids])

def main():
    bucket_name = 'flam-videoshop-assets'
    folder_name = 'flam/dev/processed_af/'
    output_file = 'output.csv'
    
    object_ids = get_s3_object_ids(bucket_name, folder_name)
    
    if object_ids:
        write_to_csv(output_file, object_ids)
        print(f'Successfully wrote the object IDs to {output_file}')
    else:
        print('No objects found in the specified folder.')

if __name__ == '__main__':
    main()
