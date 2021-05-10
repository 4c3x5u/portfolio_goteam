from rest_framework import serializers
import status

from main.serializers.user.default import UserSerializer
from main.validation.auth import authenticate_custom, authorization_error
from main.validation.custom import CustomAPIException
from main.models import User


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
        auth_user = attrs.pop('auth_user')
        auth_token = attrs.pop('auth_token')
        authenticated_user, authentication_error = \
            authenticate_custom(auth_user, auth_token)

        if authentication_error:
            raise authentication_error

        if not authenticated_user.is_admin:
            raise authorization_error

        user = attrs.get('user')
        if user.team_id != authenticated_user.team_id:
            raise authorization_error
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

