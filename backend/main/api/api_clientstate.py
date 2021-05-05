from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import User, Board
from ..validation.val_auth import not_authenticated_response, \
    not_authorized_response
from ..util import create_board
import bcrypt


@api_view(['GET'])
def client_state(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')
    board_id = request.query_params.get('boardId')

    try:
        user = User.objects.prefetch_related(
            'team',
            'board_set',
            'team__user_set',
            'team__board_set__user',
            'team__board_set__column_set',
            'team__board_set__column_set__task_set',
            'team__board_set__column_set__task_set__subtask_set'
        ).get(username=username)
    except User.DoesNotExist:
        return Response({
            'username': ErrorDetail(string='User not found.',
                                    code='not_found')
        }, 404)

    # Authenticate
    valid_token = bytes(user.username, 'utf-8') + user.password
    match = bcrypt.checkpw(valid_token, bytes(token, 'utf-8'))
    if not match:
        return not_authenticated_response

    boards = user.board_set.all()

    if not boards and user.is_admin:
        board, error_response = create_board(name='New Board',
                                             team_id=user.team_id,
                                             team_admin=user)
        if error_response:
            return error_response

        # return a list containing only the new board
        boards = [board]

    if board_id:
        try:
            board = user.board_set.get(id=board_id)
        except Board.DoesNotExist:
            return not_authorized_response
    else:
        board = user.board_set.all().first()

    if not board:
        return Response({
            'board': ErrorDetail(
                string='Please ask your team admin to add you to a board '
                       'and try again.',
                code='inactive'
            )
        }, 400)

    if board.team_id != user.team_id:
        return not_authenticated_response

    team_members = user.team.user_set.all()
    board_members = board.user.all()

    return Response({
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
        'boards': [
            {'id': board.id, 'name': board.name} for board in boards
        ],
        'activeBoard': {
            'id': board.id,
            'columns': [
                {
                    'id': column.id,
                    'order': column.order,
                    'tasks': column.task_set is not None and list(map(
                        lambda task: {
                            'id': task.id,
                            'title': task.title,
                            'description': task.description,
                            'order': task.order,
                            'user': task.user.username if task.user else '',
                            'subtasks': task.subtask_set is not None and list(map(
                                lambda subtask: {
                                    'id': subtask.id,
                                    'title': subtask.title,
                                    'order': subtask.order,
                                    'done': subtask.done
                                },
                                task.subtask_set.all()
                            ))
                        },
                        column.task_set.all()
                    ))
                } for column in board.column_set.all()
            ]
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
    }, 200)
