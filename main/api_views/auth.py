from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail

from ..serializers.userserializer import UserSerializer
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
        if user.password != request.data.get('password'):
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
