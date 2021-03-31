from rest_framework.decorators import api_view
from rest_framework.response import Response
from main.serializers.LoginSerializer import LoginSerializer


@api_view(['POST'])
def login(request):
    serializer = LoginSerializer(data=request.data)
    if serializer.is_valid():
        return Response({
            serializer.validated_data['username']: 'Login successful.'
        }, 200)
    return Response(serializer.errors)
