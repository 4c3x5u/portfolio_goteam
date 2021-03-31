from rest_framework.generics import CreateAPIView
from main.models import User
from main.serializers.RegisterSerializer import RegisterSerializer


class RegisterView(CreateAPIView):
    serializer_class = RegisterSerializer
