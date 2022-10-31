from server.main.models import User
import bcrypt
import status
from server.main.helpers.custom_api_exception import CustomAPIException


class AuthHelper:
    AUTHENTICATION_ERROR = CustomAPIException('auth',
                                              'Authentication failure.',
                                              status.HTTP_403_FORBIDDEN)

    AUTHORIZATION_ERROR = CustomAPIException('auth',
                                             'Authorization failure.',
                                             status.HTTP_401_UNAUTHORIZED)

    @classmethod
    def authenticate(cls, username, token):
        try:
            user = User.objects.get(username=username)
        except (User.DoesNotExist, ValueError):
            raise cls.AUTHENTICATION_ERROR

        try:
            tokens_match = bcrypt.checkpw(
                bytes(user.username, 'utf-8') + user.password,
                bytes(token, 'utf-8')
            )
            if not tokens_match:
                raise cls.AUTHENTICATION_ERROR
        except ValueError:
            raise cls.AUTHENTICATION_ERROR

        return user

    @classmethod
    def authorize(cls, user, team_id):
        if not user.is_admin:
            raise cls.AUTHORIZATION_ERROR
        if user.team_id != team_id:
            raise cls.AUTHORIZATION_ERROR
