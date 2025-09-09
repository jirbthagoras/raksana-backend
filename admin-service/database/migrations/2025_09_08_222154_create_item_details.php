<?php

use Illuminate\Container\Attributes\DB;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB as FacadesDB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('greenprints', function (Blueprint $table) {
            $table->id();
            $table->foreignId("item_id")->references("id")->on("items");
            $table->string("image_key");
            $table->string("title");
            $table->text("description");
            $table->string("sustainability_score");
            $table->string("estimated_time");
            $table->timestamp('created_at')->default(FacadesDB::raw('CURRENT_TIMESTAMP'));
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('greenprint');
    }
};
