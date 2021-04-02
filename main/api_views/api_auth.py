from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import bcrypt

from ..serializers.ser_user import UserSerializer
from ..models import User


@api_view(['POST'])
def register(request):
    serializer = UserSerializer(data=request.data)
    if serializer.is_valid():
        user = serializer.save()
        return Response({
            'msg': 'Registration successful.',
            'username': user.username,
        }, 201)
    else:
        return Response(serializer.errors, 400)


@api_view(['POST'])
def login(request):
    serializer = UserSerializer(data=request.data)
    if serializer.is_valid():
        try:
            user = User.objects.get(username=request.data.get('username'))
        except User.DoesNotExist:
            return Response({
                'username': ErrorDetail(string='Invalid username.',
                                        code='invalid')
            }, 400)
        pw_bytes = request.data.get('password').encode()
        if not bcrypt.checkpw(pw_bytes, bytes(user.password)):
            return Response({
                'password': ErrorDetail(string='Invalid password.',
                                        code='invalid')
            }, 400)
        return Response({
            'msg': 'Login successful.',
            'username': user.username
        }, 200)
    else:
        return Response(serializer.errors, 400)
