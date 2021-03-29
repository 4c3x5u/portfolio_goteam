from rest_framework.generics import CreateAPIView
from rest_framework.exceptions import ValidationError
from rest_framework.response import Response
from django.core.exceptions import ObjectDoesNotExist
from .serializers import UserSerializer
from .models import Team, User


class UserCreate(CreateAPIView):
    serializer_class = UserSerializer

    def create(self, request, *args, **kwargs):
        pw = request.data.get('password')
        cf = request.data.get('password_confirmation')
        ic = request.data.get('invite_code')
        if pw == cf:
            team = None
            is_admin = False
            try:
                team = Team.objects.get(invite_code=ic)
            except Team.DoesNotExist:
                team = Team.objects.create()
                is_admin = True
            finally:
                user = {'username': request.data.get('username'),
                        'password': request.data.get('password'),
                        'team': team.id,
                        'is_admin': is_admin}
                print(f'USER: {user}')
                serializer = UserSerializer(data=user)
                if serializer.is_valid():
                    serializer.save()
                    return Response(user)
                else:
                    is_admin and team.delete()
                    return Response(serializer.errors )
        return Response({'password': "confirmation doesn't match password"})


# @api_view(['POST'])
# def team(request):
#     return HttpResponse(Team.objects.create())
