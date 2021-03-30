from rest_framework.generics import CreateAPIView
from rest_framework.response import Response
from main.serializers import UserSerializer
from main.models import Team


class Register(CreateAPIView):
    serializer_class = UserSerializer

    def create(self, request, *args, **kwargs):
        serializer = UserSerializer(data=request.data)
        if serializer.is_valid():
            try:
                serializer.save()
                return super().create(request, *args, **kwargs)
            except Team.DoesNotExist:
                return Response({'invite_code': "Team not found."}, 404)
        else:
            return Response(serializer.errors, 400)


