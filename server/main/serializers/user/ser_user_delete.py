from rest_framework import serializers
import status

from server.main.serializers.user.ser_user import UserSerializer
from server.main.helpers.auth_helper import AuthHelper
from server.main.helpers.custom_api_exception import CustomAPIException
from server.main.models import User


class DeleteUserSerializer(UserSerializer):
    """
    Used only for adding/removing a user to/from a board
    """
    user = serializers.PrimaryKeyRelatedField(
        queryset=User.objects.all(),
        error_messages={
            'null': 'Username cannot be null.',
            'does_not_exist': 'User does not exist.'
        }
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = UserSerializer.Meta.model
        fields = 'user', 'auth_user', 'auth_token',

    def validate(self, attrs):
        authenticated_user = AuthHelper.authenticate(attrs.pop('auth_user'),
                                                     attrs.pop('auth_token'))
        user = attrs.get('user')
        AuthHelper.authorize(authenticated_user, user.team_id)
        if user.is_admin:
            raise CustomAPIException(
                'username',
                'Admins cannot be deleted from their teams.',
                status.HTTP_403_FORBIDDEN
            )
        return user

    def delete(self):
        self.instance = {'username': self.validated_data.username}
        return self.validated_data.delete()

    def to_representation(self, instance):
        return {
            'msg': 'Member has been deleted successfully.',
            'username': instance.get('username'),
        }

