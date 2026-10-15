const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");

const PROTO_PATH = "../user.proto";

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const userPackage = protoDescriptor.user;

const client = new userPackage.UserService(
    "localhost:50051",
    grpc.credentials.createInsecure()
);

client.GetUserProfile({ user_id: 1 }, (error, response) => {
    if (error) {
        console.error("Error:", error.message);
        client.close();
        return;
    }

    console.log("Nama:", response.name);
    console.log("Email:", response.email);

    client.close();
});