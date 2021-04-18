from rest_framework.test import APITestCase
from ..models import Team
from ..util import new_admin, new_member


class DeleteUserTests(APITestCase):
    def setUp(self):
        self.team = Team.objects.create()
        self.admin = new_admin(self.team)
        self.member = new_member(self.team)

    def delete_user(self, username, auth_user, auth_token):
        return self.client.delete(f'/users/?username={username}',
                                  HTTP_AUTH_USER=auth_user,
                                  HTTP_AUTH_TOKEN=auth_token)

    def test_success(self):
        response = self.delete_user(self.member['username'],
                                    self.admin['username'],
                                    self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Member is removed successfuly.',
            'username': self.member['username']
        })
