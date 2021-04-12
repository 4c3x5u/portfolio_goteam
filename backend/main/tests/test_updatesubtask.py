from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask, Task, Column, Board, Team, User


class UpdateSubtaskTests(APITestCase):
    def setUp(self):
        self.url = '/subtasks/'
        team = Team.objects.create()
        self.subtask = Subtask.objects.create(
            title='Some Task Title',
            order=0,
            task=Task.objects.create(
                title="Some Subtask Title",
                order=0,
                column=Column.objects.create(
                    order=0,
                    board=Board.objects.create(team=team)
                )
            )
        )
        self.admin = User.objects.create(
            username='teamadmin',
            password=b'$2b$12$lrkDnrwXSBU.YJvdzbpAWOd9GhwHJGVYafRXTHct2gm3akPJ'
                     b'gB5Zq',
            is_admin=True,
            team=team
        )
        self.member = User.objects.create(
            username='teammember',
            password=b'$2b$12$RonFQ1/18JiCN8yFeBaxKOsVbxLdcehlZ4e0r9gtZbARqEVU'
                     b'HHEoK',
            is_admin=False,
            team=team
        )
        self.admin_token = '$2b$12$TVdxI.a.ZlOkhH1/mZQ/IOHmKxklQJWiB0n6ZSg2R' \
                           'JJO17thjLOdy'
        self.member_token = '$2b$12$xnIJLzpgNV12O80XsakMjezCFqwIphdBy5ziJ9Eb' \
                            '9stnDZze19Ude'
        self.forbidden_response = {
            'auth': ErrorDetail(string="Authentication failure.",
                                code='not_authenticated')
        }

    def help_test_success(self, request_data):
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {'msg': 'Subtask update successful.',
                                         'id': self.subtask.id})
        return Subtask.objects.get(id=self.subtask.id)

    def help_test_failure(self):
        subtask = Subtask.objects.get(id=self.subtask.id)
        self.assertEqual(subtask.title, self.subtask.title)
        self.assertEqual(subtask.done, self.subtask.done)

    def test_title_success(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        subtask = self.help_test_success(request_data)
        self.assertEqual(subtask.title, request_data.get('data').get('title'))

    def test_done_success(self):
        request_data = {'id': self.subtask.id, 'data': {'done': True}}
        subtask = self.help_test_success(request_data)
        self.assertEqual(subtask.done, request_data.get('data').get('done'))

    def test_order_success(self):
        request_data = {'id': self.subtask.id, 'data': {'order': 10}}
        subtask = self.help_test_success(request_data)
        self.assertEqual(subtask.order, request_data.get('data').get('order'))

    def test_id_blank(self):
        request_data = {'id': '', 'data': {'title': 'New task title.'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'id': ErrorDetail(string='Subtask ID cannot be empty.',
                              code='blank')
        })
        self.help_test_failure()

    def test_data_blank(self):
        request_data = {'id': self.subtask.id, 'data': ''}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data': ErrorDetail(string='Data cannot be empty.', code='blank')
        })
        self.help_test_failure()

    def test_title_blank(self):
        request_data = {'id': self.subtask.id, 'data': {'title': ''}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.title': ErrorDetail(string='Title cannot be empty.',
                                      code='blank')
        })
        self.help_test_failure()

    def test_done_blank(self):
        request_data = {'id': self.subtask.id, 'data': {'done': ''}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.done': ErrorDetail(string='Done cannot be empty.',
                                     code='blank')
        })
        self.help_test_failure()

    def test_order_blank(self):
        request = {'id': self.subtask.id, 'data': {'order': ''}}
        response = self.client.patch(self.url,
                                     request,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.order': ErrorDetail(string='Order cannot be empty.',
                                      code='blank')
        })
        self.help_test_failure()

    def test_auth_token_empty(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN='')
        print(f'response: {response.data}')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)

    def test_auth_token_invalid(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosi')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)

    def test_auth_user_blank(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER='',
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)

    def test_auth_user_invalid(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER='invalidio',
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)

    def test_unauthorized(self):
        request_data = {'id': self.subtask.id,
                        'data': {'title': 'New Task Title'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.member.username,
                                     HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='The user is not an admin.',
                                code='not_authorized')
        })
