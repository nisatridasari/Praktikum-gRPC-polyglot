import grpc
import user_pb2
import user_pb2_grpc


def main():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)

        request = user_pb2.UserRequest(user_id=1)
        response = stub.GetUserProfile(request)

        print("Nama:", response.name)
        print("Email:", response.email)


if __name__ == "__main__":
    main()