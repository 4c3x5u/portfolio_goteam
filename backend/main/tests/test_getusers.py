from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Team, User
from ..util import not_authenticated_response, new_member


class GetUsersTests(APITestCase):
    endpoint = '/users/?team_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.users = [
            User.objects.create(
                username=f'User #{i}',
                password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTI'
                         b'MeeJpC6m',
                is_admin=False,
                team=self.team
            ) for i in range(0, 3)
        ]
        self.token = '$2b$12$WLmxQnf9kbDoW/8jA6kfIO9TfchCiGphBpckS2oy755wtdT' \
                     'aIQsoq'

    def test_success(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, list(map(
            lambda user: {'username': user.username},
            self.users
        )))

    def test_team_id_blank(self):
        response = self.client.get(f'{self.endpoint}',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })

    def test_team_id_invalid(self):
        response = self.client.get(f'{self.endpoint}qwerty',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID must be a number.',
                                   code='invalid')
        })

    def test_team_not_found(self):
        response = self.client.get(f'{self.endpoint}12412312',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        })


    def test_auth_user_empty(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER='',
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER='invalidusername',
                                   HTTP_AUTH_TOKEN=self.token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_empty(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.users[0].username,
                                   HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosia')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

