from rest_framework.generics import CreateAPIView
from main.serializers.RegisterSerializer import RegisterSerializer


class RegisterAPIView(CreateAPIView):
    serializer_class = RegisterSerializer
