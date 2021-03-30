from rest_framework.generics import CreateAPIView
from rest_framework.response import Response
from main.serializers import UserSerializer
from main.models import Team


class Register(CreateAPIView):
    serializer_class = UserSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.serializer_class(data=request.data)
        if serializer.is_valid():
            return super().create(request, *args, **kwargs)
        else:
            return Response(serializer.errors, 400)


