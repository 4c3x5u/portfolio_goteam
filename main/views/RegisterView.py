from rest_framework.generics import CreateAPIView
from rest_framework.response import Response
from main.serializers.RegisterSerializer import RegisterSerializer


class RegisterView(CreateAPIView):
    serializer_class = RegisterSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.serializer_class(data=request.data)
        if serializer.is_valid():
            return super().create(request, *args, **kwargs)
        else:
            return Response(serializer.errors, 400)
