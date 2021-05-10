from rest_framework.decorators import api_view
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
import status

from ..models import Board
from ..serializers.updateuserserializer import UpdateUserSerializer
from ..validation.val_auth import \
    authenticate, authorize, not_authorized_response
from ..validation.val_board import validate_board_id
from ..validation.val_user import validate_username, validate_is_active


class Users(APIView):
    @staticmethod
    def patch(request):
        """
        Used only for adding/removing a user to/from a board
        """
        serializer = UpdateUserSerializer(data={
            'username': request.query_params.get('username'),
            'board_id': request.data.get('board_id') or None,
            'is_active': request.data.get('is_active'),
            'auth_user': request.META.get('HTTP_AUTH_USER'),
            'auth_token': request.META.get('HTTP_AUTH_TOKEN')
        })
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status.HTTP_200_OK)
        return Response(serializer.errors, status.HTTP_400_BAD_REQUEST)


@api_view(['PATCH', 'DELETE'])
def users(request):
    auth_user = request.META.get('HTTP_AUTH_USER')
    auth_token = request.META.get('HTTP_AUTH_TOKEN')

    auth_user, authentication_response = authenticate(auth_user, auth_token)
    if authentication_response:
        return authentication_response

    if request.method == 'DELETE':
        authorization_response = authorize(auth_user.username)
        if authorization_response:
            return authorization_response

        username = request.query_params.get('username')
        user, validation_response = validate_username(username)
        if validation_response:
            return validation_response
        if user.team_id != auth_user.team_id:
            return not_authorized_response

        # this is not authorization. it checks whether the user that is up for
        # deletion is admin
        if user.is_admin:
            return Response({
                'username': ErrorDetail(
                    string='Team leaders cannot be deleted from their teams.',
                    code='forbidden'
                )
            }, 403)

        user.delete()

        return Response({
            'msg': 'Member has been deleted successfully.',
        }, 200)
