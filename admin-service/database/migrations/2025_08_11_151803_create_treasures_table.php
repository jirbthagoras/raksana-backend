<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('treasures', function (Blueprint $table) {
            $table->id();
            $table->string("name");
            $table->bigInteger("point_gain");
            $table->string("code_id", 255);
            $table->foreign("code_id")->references("id")->on("codes");
            $table->boolean("claimed")->default(false);
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('treasures');
    }
};
