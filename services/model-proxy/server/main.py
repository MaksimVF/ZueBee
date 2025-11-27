



import asyncio
import grpc
import model_proxy_pb2
import model_proxy_pb2_grpc
from litellm_backend import generate_text

class ModelProxyServicer(model_proxy_pb2_grpc.ModelProxyServicer):
    async def Generate(self, request, context):
        try:
            async for part in generate_text(
                model=request.model,
                prompt=request.prompt,
            ):
                yield model_proxy_pb2.GenerateResponse(chunk=part)
        except Exception as e:
            yield model_proxy_pb2.GenerateResponse(error=str(e), done=True)

async def serve():
    server = grpc.aio.server()
    model_proxy_pb2_grpc.add_ModelProxyServicer_to_server(
        ModelProxyServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    await server.start()
    await server.wait_for_termination()

if __name__ == "__main__":
    asyncio.run(serve())



