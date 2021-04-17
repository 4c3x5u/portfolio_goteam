from rest_framework.test import APITestCase
from ..models import Team, User
from ..util import new_member


class GetUsersTests(APITestCase):
    endpoint = '/users/team_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.member = new_member(self.team)
        self.users = [
            User.objects.create(
                username=f'User #{i}',
                password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTI'
                         b'MeeJpC6me',
                is_admin=False,
                team=self.team
            ) for i in range(0, 3)
        ]

    def test_success(self):
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN=self.member.token)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(
            response.data,
            list(map(
                lambda user: {'username': user.username},
                self.users
            ))
        )
