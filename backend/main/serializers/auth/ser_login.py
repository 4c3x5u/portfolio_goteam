import bcrypt
import status

from main.models import User
from main.serializers.user.ser_user import UserSerializer
from main.validation.val_custom import CustomAPIException


class LoginSerializer(UserSerializer):
    def validate(self, attrs):
        try:
            user = User.objects.get(username=attrs.get('username'))
        except User.DoesNotExist:
            raise CustomAPIException('username',
                                     'Invalid username.',
                                     status.HTTP_400_BAD_REQUEST)

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
