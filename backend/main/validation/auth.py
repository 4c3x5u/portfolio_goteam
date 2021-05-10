from main.models import User
from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response
import bcrypt
import status
from .custom import CustomAPIException


# TODO: delete once you moved on to the customvalidation approach
not_authenticated_response = Response({
    'auth': ErrorDetail(string="Authentication failure.",
                        code='not_authenticated')
}, 403)


# TODO: delete once you moved on to the customvalidation approach
def authenticate(username, token):
    try:
        user = User.objects.get(username=username)
    except (User.DoesNotExist, ValueError):
        return None, not_authenticated_response

    try:
        tokens_match = bcrypt.checkpw(
            bytes(user.username, 'utf-8') + user.password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return None, not_authenticated_response
    except ValueError:
        return None, not_authenticated_response

    return user, None

# TODO: delete once you moved on to the customvalidation approach
not_authorized_response = Response({
    'auth': ErrorDetail(string='Authorization failure.',
                        code='not_authorized')
}, 403)


# TODO: delete once you moved on to the customvalidation approach
def authorize(username):
    try:
        user = User.objects.get(username=username)
        if not user.is_admin:
            return not_authorized_response
    except User.DoesNotExist:
        return not_authorized_response


authentication_error = CustomAPIException('auth',
                                          'Authentication failure.',
                                          status.HTTP_403_FORBIDDEN)


def authenticate_custom(username, token):
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


# TODO: rename as authorize once you moved on to the customvalidation approach
def authorize_custom(username):
    try:
        user = User.objects.get(username=username)
        if not user.is_admin:
            return authorization_error
    except User.DoesNotExist:
        return authorization_error
