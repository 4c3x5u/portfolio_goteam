from main.models import Team, User
from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response
import bcrypt


def new_admin(team):
    user = User.objects.create(
        username='teamadmin',
        password=b'$2b$12$lrkDnrwXSBU.YJvdzbpAWOd9GhwHJGVYafRXTHct2gm3akPJgB5Z'
                 b'q',
        is_admin=True,
        team=team
    )
    token = '$2b$12$TVdxI.a.ZlOkhH1/mZQ/IOHmKxklQJWiB0n6ZSg2RJJO17thjLOdy'
    return {'username': user.username,
            'password': user.password,
            'is_admin': user.is_admin,
            'team': user.team,
            'token': token}


def new_member(team):
    user = User.objects.create(
        username='teammember',
        password=b'$2b$12$RonFQ1/18JiCN8yFeBaxKOsVbxLdcehlZ4e0r9gtZbARqEVUHHEo'
                 b'K',
        is_admin=False,
        team=team
    )
    token = '$2b$12$xnIJLzpgNV12O80XsakMjezCFqwIphdBy5ziJ9Eb9stnDZze19Ude'
    return {'username': user.username,
            'password': user.password,
            'is_admin': user.is_admin,
            'team': user.team,
            'token': token}


forbidden_response = {'auth': ErrorDetail(string="Authentication failure.",
                                          code='not_authenticated')}


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


