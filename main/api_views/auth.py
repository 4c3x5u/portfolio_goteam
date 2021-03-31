from rest_framework.decorators import api_view
from rest_framework.response import Response
from main.serializers.register_serializer import RegisterSerializer
from main.serializers.login_serializer import LoginSerializer


@api_view(['POST'])
def register(request):
    serializer = RegisterSerializer(data=request.data)
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
    serializer = LoginSerializer(data=request.data)
    if serializer.is_valid():
        return Response({
            'msg': 'Login successful.',
            'username': serializer.validated_data['username'],
        }, 200)
    return Response(serializer.errors, 400)
