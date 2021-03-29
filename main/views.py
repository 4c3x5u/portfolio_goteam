from rest_framework.decorators import api_view
from rest_framework.response import Response
from uuid import UUID
from .serializers import UserSerializer
from .models import Team


@api_view(['POST'])
def register(request):
    pw = request.data.get('password')
    if not pw:
        return Response({'password': 'Password cannot be empty.'}, 400)
    cf = request.data.get('password_confirmation')
    if not cf:
        return Response({
            'password_confirmation': 'Password confirmation cannot be empty.'
        }, 400)
    if pw == cf:
        ic = request.data.get('invite_code')
        if ic:
            try:
                UUID(ic)
            except ValueError:
                return Response({'invite_code': 'Invalid invite code.'}, 400)
            try:
                team = Team.objects.get(invite_code=ic)
                is_admin = False
            except Team.DoesNotExist:
                return Response({'invite_code': "team not found"}, 404)
        else:
            team = Team.objects.create()
            is_admin = True
        new_user = {'username': request.data.get('username'),
                    'password': request.data.get('password'),
                    'team': team.id,
                    'is_admin': is_admin}
        serializer = UserSerializer(data=new_user)
        if serializer.is_valid():
            serializer.save()
            return Response(new_user, 201)
        else:
            is_admin and team.delete()
            return Response(serializer.errors, 400)
    return Response({
        'password_confirmation': "Confirmation does not match the password."
    }, 400)
