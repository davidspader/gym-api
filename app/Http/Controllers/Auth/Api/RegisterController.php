<?php

namespace App\Http\Controllers\Auth\Api;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Http\Request;

class RegisterController extends Controller
{
    public function register(Request $request, User $user)
    {
        if(!$request->has('name', 'email', 'password')){
            return response()->json( [
                'error' => [
                    'message' => 'Error to create a new user'
                ]
            ], 401);
        }

        $userData = $request->only('name', 'email', 'password');
        $userData['password'] = bcrypt($userData['password']);

        if(!$user = $user->create($userData)) {
            return response()->json( [
                'error' => [
                    'message' => 'Error to create a new user'
                ]
            ], 401);
        }

        return response()->json([
           'data' => [
               'user' => $user,
           ]
        ]);
    }
}

