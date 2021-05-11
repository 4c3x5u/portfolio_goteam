from main.models import User
import bcrypt
import status
from .val_custom import CustomAPIException

authentication_error = CustomAPIException('auth',
                                          'Authentication failure.',
                                          status.HTTP_403_FORBIDDEN)

authorization_error = CustomAPIException('auth',
                                         'Authorization failure.',
                                         status.HTTP_401_UNAUTHORIZED)


def authenticate(username, token):
    try:
        user = User.objects.get(username=username)
    except (User.DoesNotExist, ValueError):
        raise authentication_error

    try:
        tokens_match = bcrypt.checkpw(
            bytes(user.username, 'utf-8') + user.password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            raise authentication_error
    except ValueError:
        raise authentication_error

    return user


def authorize(user, team_id):
    if not user.is_admin:
        raise authorization_error
    if user.team_id != team_id:
        raise authorization_error
