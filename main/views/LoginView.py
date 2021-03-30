from rest_framework.generics import GenericAPIView
from rest_framework.response import Response
from main.models import User


class LoginView(GenericAPIView):
    queryset = User.objects.all()

    def post(self, request, *args, **kwargs):
        if not request.data['username']:
            return Response({
                'username': 'Username cannot be empty.'
            }, 400)
        if not request.data['password']:
            return Response({
                'password': 'Password cannot be empty.'
            }, 400)
        user = self.queryset.get(username=request.data['username'])
        if not user:
            return Response({
                'password': 'Username not found.'
            }, 404)
        if user.password != request.data['password']:
            return Response('Incorrect password.', 400)
        else:
            return Response('success', 200)
