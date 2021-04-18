from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import User, Board, Team
from ..util import new_admin


class PostUsersTests(APITestCase):
    def setUp(self):
        team = Team.objects.create()
        self.admin = new_admin(team)
        self.user = User.objects.create(
            username='Some User',
            password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTIMeeJ'
                     b'pC6m',
            is_admin=False,
            team=team
        )
        self.board = Board.objects.create(name='Some Board', team=team)
        self.board.user.add(self.user)
        self.username = self.user.username
        self.token = '$2b$12$WLmxQnf9kbDoW/8jA6kfIO9TfchCiGphBpckS2oy755wtdT' \
                     'aIQsoq'

    def postUser(self, user_data, auth_user, auth_token):
        return self.client.post(f'/users/',
                                user_data,
                                HTTP_AUTH_USER=auth_user,
                                HTTP_AUTH_TOKEN=auth_token)

    def test_success(self):
        response = self.postUser({
            'username': self.user.username,
            'board_id': self.board.id,
            'is_active': False
        }, self.admin['username'], self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': f'{self.user.username} is removed from {self.board.name}.'
        })
        self.assertEqual(
            len(self.board.user.filter(username=self.user.username)),
            0
        )

    def test_username_blank(self):
        response = self.postUser({
            'username': '',
            'board_id': self.board.id,
            'is_active': False
        }, self.admin['username'], self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Username cannot be empty.',
                                    code='blank')
        })
