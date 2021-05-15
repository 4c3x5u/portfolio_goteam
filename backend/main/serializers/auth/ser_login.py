from rest_framework import serializers
import bcrypt
import status

from ...models import User
from ..user.ser_user import UserSerializer
from main.helpers.custom_api_exception import CustomAPIException


class LoginSerializer(UserSerializer):
    username = serializers.PrimaryKeyRelatedField(
        queryset=User.objects.all(),
        error_messages={'does_not_exist': 'Invalid username.',
                        'null': 'Username cannot be null.'}
    )

    def validate(self, attrs):
        user = attrs.get('username')
        pw_bytes = bytes(attrs.get('password'), 'utf-8')
        if not bcrypt.checkpw(pw_bytes, bytes(user.password)):
            raise CustomAPIException('password',
                                     'Invalid password.',
                                     status.HTTP_400_BAD_REQUEST)
        return user

    def to_representation(self, instance):
        return {
            'msg': 'Login successful.',
            'username': instance.username,
            'token': bcrypt.hashpw(
                bytes(instance.username, 'utf-8') + instance.password,
                bcrypt.gensalt()
            ).decode('utf-8'),
            'teamId': instance.team_id,
            'isAdmin': instance.is_admin,
        }
