from main.models import User
import bcrypt
import status
from .val_custom import CustomAPIException


authentication_error = CustomAPIException('auth',
                                          'Authentication failure.',
                                          status.HTTP_403_FORBIDDEN)


def authenticate(username, token):
    try:
        user = User.objects.get(username=username)
    except (User.DoesNotExist, ValueError):
        return None, authentication_error

    try:
        tokens_match = bcrypt.checkpw(
            bytes(user.username, 'utf-8') + user.password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return None, authentication_error
    except ValueError:
        return None, authentication_error

    return user, None


authorization_error = CustomAPIException('auth',
                                         'Authorization failure.',
                                         status.HTTP_401_UNAUTHORIZED)


def authorize(username):
    try:
        user = User.objects.get(username=username)
        if not user.is_admin:
            return authorization_error
    except User.DoesNotExist:
        return authorization_error
