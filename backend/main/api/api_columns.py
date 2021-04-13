from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..util import authenticate, authorize
from ..models import Column
from ..serializers.ser_column import ColumnSerializer


@api_view(['GET'])
def columns(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    authorization_response = authorize(username)
    if authorization_response:
        return authorization_response

    board_id = request.query_params.get('board_id')

    if not board_id:
        return Response({
            'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                    code='blank')
        }, 400)

    try:
        int(board_id)
    except ValueError:
        return Response({
            'board_id': ErrorDetail(string='Board ID must be a number.',
                                    code='invalid')
        }, 400)

    board_columns = Column.objects.filter(board_id=board_id)
    serializer = ColumnSerializer(board_columns, many=True)

    if board_columns:
        return Response({
            'columns': list(
                map(
                    lambda column: {
                        'id': column['id'],
                        'order': column['order']
                    },
                    serializer.data
                )
            )
        }, 200)
