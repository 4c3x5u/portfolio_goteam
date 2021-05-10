from rest_framework.views import APIView
from rest_framework.response import Response
import status

from main.serializers.auth.loginserializer import LoginSerializer


class LoginAPIView(APIView):
    @staticmethod
    def post(request):
        serializer = LoginSerializer(data=request.data)
        if not serializer.is_valid():
            return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
        return Response(serializer.data, status.HTTP_200_OK)