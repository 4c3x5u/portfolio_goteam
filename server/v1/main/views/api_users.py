from rest_framework.views import APIView
from rest_framework.response import Response
import status

from main.serializers.user.ser_user_update import UpdateUserSerializer
from main.serializers.user.ser_user_delete import DeleteUserSerializer


class UsersApiView(APIView):
    @staticmethod
    def patch(request):
        """
        Used only for adding/removing a user to/from a board
        """
        serializer = UpdateUserSerializer(data={
            'username': request.query_params.get('username'),
            'board': request.data.get('board_id') or None,
            'is_active': request.data.get('is_active'),
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def delete(request):
        serializer = DeleteUserSerializer(data={
            'user': request.query_params.get('username'),
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.delete()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)
