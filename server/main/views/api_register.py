from rest_framework.views import APIView
from rest_framework.response import Response
import status

from server.main.serializers.auth.ser_register import RegisterSerializer


class RegisterApiView(APIView):
    @staticmethod
    def post(request):
        invite_code = request.query_params.get('invite_code')
        serializer = RegisterSerializer(data={
            'username': request.data.get('username'),
            'password': request.data.get('password'),
            'password_confirmation': request.data.get('password_confirmation'),
            'invite_code': invite_code
        } if invite_code else request.data)
        if not serializer.is_valid():
            return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
        serializer.save()
        return Response(serializer.data, status.HTTP_201_CREATED)
