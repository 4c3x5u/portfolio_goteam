from rest_framework.decorators import api_view
from rest_framework.response import Response
from main.serializers.register_serializer import RegisterSerializer


@api_view(['POST'])
def register(request):
    serializer = RegisterSerializer(data=request.data)
    if serializer.is_valid():
        user = serializer.save()
        return Response({
            'msg': 'Login successful.',
            'username': user.username
        }, 201)
    else:
        return Response(serializer.errors, 400)
