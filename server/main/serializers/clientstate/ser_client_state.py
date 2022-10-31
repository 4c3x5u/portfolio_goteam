from rest_framework import serializers
from server.main.models import User, Board
from server.main.helpers.auth_helper import AuthHelper
from server.main.helpers.custom_api_exception import CustomAPIException
import bcrypt
import status


class ClientStateSerializer(serializers.Serializer):
    user = serializers.PrimaryKeyRelatedField(
        queryset=User.objects.prefetch_related(
            'team',
            'board_set',
            'team__user_set',
            'team__board_set__user',
            'team__board_set__column_set',
            'team__board_set__column_set__task_set',
            'team__board_set__column_set__task_set__subtask_set'
        ).all(),
    )
    auth_token = serializers.CharField()
    board_id = serializers.IntegerField(default=-1)

    def create(self, validated_data):
        pass

    def update(self, instance, validated_data):
        pass

    def validate(self, attrs):
        user = attrs.get('user')

        valid_token = bytes(user.username, 'utf-8') + user.password
        match = bcrypt.checkpw(
            valid_token,
            bytes(attrs.get('auth_token'), 'utf-8')
        )
        if not match:
            raise AuthHelper.AUTHENTICATION_ERROR

        board_id = attrs.get('board_id')
        if board_id:
            try:
                board = user.board_set.get(id=board_id)
            except Board.DoesNotExist:
                raise AuthHelper.AUTHORIZATION_ERROR
        else:
            board = user.board_set.all().first()

        if not board:
            err_detail = ('Please ask your team admin to add you to a board '
                          'and try again.')
            raise CustomAPIException('board',
                                     err_detail,
                                     status.HTTP_400_BAD_REQUEST)

        if board.team_id != user.team_id:
            raise AuthHelper.AUTHORIZATION_ERROR

        return {'user': user, 'board': board}

    def to_representation(self, instance):
        user = instance.get('user')
        board = instance.get('board')

        team_members = user.team.user_set.all()
        board_members = board.user.all()
        boards = user.board_set.all()

        return {
            'user': {
                'username': user.username,
                'teamId': user.team_id,
                'isAdmin': user.is_admin,
                'isAuthenticated': True
            },
            'team': user.is_admin and {
                'id': user.team.id,
                'inviteCode': user.team.invite_code
            },
            'boards': [{
                'id': board.id, 'name': board.name
            } for board in boards],
            'activeBoard': {
                'id': board.id,
                'columns': [{
                    'id': column.id,
                    'order': column.order,
                    'tasks': column.task_set is not None and [{
                        'id': task.id,
                        'title': task.title,
                        'description': task.description,
                        'order': task.order,
                        'user': task.user.username if task.user else '',
                        'subtasks': task.subtask_set is not None and [{
                            'id': subtask.id,
                            'title': subtask.title,
                            'order': subtask.order,
                            'done': subtask.done
                        } for subtask in task.subtask_set.all()]
                    } for task in column.task_set.all()]
                } for column in board.column_set.all()]
            },
            'members': [
                {
                    'username': member.username,
                    'isActive': member in board_members,
                    'isAdmin': member.is_admin
                } for member in sorted(
                    team_members,
                    key=lambda member: not member.is_admin
                )
            ]
        }


