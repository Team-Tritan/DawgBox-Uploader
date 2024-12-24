import { Client } from "minio";
import { createReadStream } from "fs";
import { extname } from "path";

let config = {
  endpoint: "s3.mci.bb.tritan.host",
  publicURL: "https://s3.tritan.gg",
  bucketName: "uploads",
};

async function main(filePath: string) {
  if (!filePath) {
    console.error("No file path provided.\nUsage: ./executable <file-path>");
    process.exit(1);
  }

  let file;
  try {
    file = createReadStream(filePath);
  } catch (err) {
    console.error(`Could not open file: ${filePath}`);
    process.exit(1);
  }

  const fileExtension = extname(filePath);
  const fileName = generateRandomString(10);
  const newFileName = fileName + fileExtension;

  const minioClient = new Client({
    endPoint: config.endpoint,
    accessKey: "",
    secretKey: "",
    useSSL: false,
  });

  try {
    await minioClient.putObject(config.bucketName, newFileName, file);

    const url = `${config.publicURL}/${config.bucketName}/${newFileName}`;
    console.log(`File uploaded successfully.\nURL: ${url}`);
  } catch (err) {
    console.error(`Failed to upload file: ${err}`);
    process.exit(1);
  }
}

function generateRandomString(length: number) {
  const charset =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += charset[Math.floor(Math.random() * charset.length)];
  }
  return result;
}

main(process.argv[2]);
