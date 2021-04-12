from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import User, Team
import bcrypt


def authenticate(username, token):
    not_authenticated_response = Response({
        'auth': ErrorDetail(string="Authentication failure.",
                            code='not_authenticated')
    }, 403)

    try:
        user = User.objects.get(username=username)
    except (User.DoesNotExist, ValueError):
        return not_authenticated_response

    try:
        tokens_match = bcrypt.checkpw(
            bytes(user.username, 'utf-8') + user.password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return not_authenticated_response
    except ValueError:
        return not_authenticated_response


def authorize(username):
    not_authorized_response = Response({
        'auth': ErrorDetail(
            string='The user is not an admin.',
            code='not_authorized'
        )
    }, 403)

    try:
        user = User.objects.get(username=username)
        if not user.is_admin:
            return not_authorized_response
    except User.DoesNotExist:
        return not_authorized_response



def validate_team_id(team_id):
    if not team_id:
        return Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        }, 400)
    try:
        Team.objects.get(id=team_id)
    except Team.DoesNotExist:
        return Response({
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        }, 404)


